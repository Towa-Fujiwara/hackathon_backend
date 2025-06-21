package usecase

import (
    "context"
    "encoding/json"
    "fmt"
    "regexp"
    "strings"
    "hackathon/dao"
    aiplatform "cloud.google.com/go/aiplatform/apiv1"
    "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
    "google.golang.org/api/option"
    "log"
)

type GeminiUsecase struct {
    postDao dao.PostDao
    client  *aiplatform.PredictionClient
    projectID string
    location  string
    modelID   string
}

func NewGeminiUsecase(postDao dao.PostDao, projectID, location, engineID string) (*GeminiUsecase, error) {
    ctx := context.Background()

    // Vertex AIクライアントの初期化
    client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)))
    if err != nil {
        return nil, fmt.Errorf("failed to create PredictionClient: %v", err)
    }

    return &GeminiUsecase{
        postDao: postDao,
        client:  client,
        projectID: projectID,
        location:  location,
        modelID:   engineID,
    }, nil
}

type UserSummary struct {
    UserId    string `json:"userId"`
    UserName  string `json:"userName"`
    Summary   string `json:"summary"`
    Interests []string `json:"interests"`
    Personality string `json:"personality"`
}

type GeminiResponse struct {
    Summary    string   `json:"summary"`
    Interests  []string `json:"interests"`
    Personality string  `json:"personality"`
}

func (g *GeminiUsecase) GenerateUserSummary(ctx context.Context, userId string) (*UserSummary, error) {
    log.Printf("GeminiUsecase: GenerateUserSummaryで受信したuserId: %s", userId)

    posts, err := g.postDao.FindAllByUserId(userId)
    if err != nil {
        log.Printf("GeminiUsecase: g.postDao.FindAllByUserId(%s) エラー: %v", userId, err)
        return nil, fmt.Errorf("failed to get user posts: %v", err)
    }
    log.Printf("GeminiUsecase: g.postDao.FindAllByUserId(%s) で取得した投稿数: %d", userId, len(posts))

    if len(posts) == 0 {
        return &UserSummary{
            UserId:    userId,
            UserName:  "Unknown",
            Summary:   "まだ投稿がありません",
            Interests: []string{},
            Personality: "投稿が少ないため、性格を分析できません",
        }, nil
    }

    var postTexts []string
    for _, post := range posts {
        if post.Text != "" {
            postTexts = append(postTexts, post.Text)
        }
    }

    if len(postTexts) == 0 {
        return &UserSummary{
            UserId:    userId,
            UserName:  posts[0].UserName,
            Summary:   "テキスト投稿がありません",
            Interests: []string{},
            Personality: "テキスト投稿が少ないため、性格を分析できません",
        }, nil
    }

    allText := strings.Join(postTexts, "\n\n")

    prompt := fmt.Sprintf(`
以下のユーザーの投稿を分析して、以下の形式でJSONを返してください：

{
  "summary": "ユーザーの投稿内容を200文字以内で要約",
  "interests": ["興味・関心のある分野を配列で"],
  "personality": "投稿から読み取れる性格・特徴を100文字以内で"
}

ユーザーの投稿：
%s

JSONのみを返してください。`, allText)

    // プロンプトの準備
    instance := &aiplatformpb.Value{
        Kind: &aiplatformpb.Value_StructValue{
            StructValue: &aiplatformpb.Struct{
                Fields: map[string]*aiplatformpb.Value{
                    "prompt": {Kind: &aiplatformpb.Value_StringValue{StringValue: prompt}},
                },
            },
        },
    }

    // パラメータの設定
    parameters := &aiplatformpb.Value{
        Kind: &aiplatformpb.Value_StructValue{
            StructValue: &aiplatformpb.Struct{
                Fields: map[string]*aiplatformpb.Value{
                    "temperature": {Kind: &aiplatformpb.Value_NumberValue{NumberValue: 0.2}},
                    "maxOutputTokens": {Kind: &aiplatformpb.Value_NumberValue{NumberValue: 1000}},
                },
            },
        },
    }

    // 予測リクエストの構築
    req := &aiplatformpb.PredictRequest{
        Endpoint:   fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s", g.projectID, g.location, g.modelID),
        Instances:  []*aiplatformpb.Value{instance},
        Parameters: parameters,
    }

    // 予測の実行
    resp, err := g.client.Predict(ctx, req)
    if err != nil {
        log.Printf("GeminiUsecase: Gemini API Predict エラー: %v", err)
        return nil, fmt.Errorf("failed to predict: %v", err)
    }

    if len(resp.GetPredictions()) == 0 {
        log.Print("GeminiUsecase: Gemini APIからのレスポンス候補がない")
        return nil, fmt.Errorf("no response from Gemini API")
    }

    // レスポンスからテキストを抽出
    prediction := resp.GetPredictions()[0]
    var summaryText string
    
    if contentValue, ok := prediction.GetStructValue().GetFields()["content"]; ok {
        if textValue, ok := contentValue.GetStringValueOk(); ok {
            summaryText = textValue
        } else {
            // contentフィールドが存在するが、StringValueではない場合
            log.Printf("GeminiUsecase: contentフィールドがStringValueではありません: %+v", contentValue)
            return nil, fmt.Errorf("unexpected content type in response")
        }
    } else {
        // contentフィールドが存在しない場合、prediction全体を文字列として扱う
        log.Printf("GeminiUsecase: contentフィールドが見つかりません。prediction全体: %+v", prediction)
        return nil, fmt.Errorf("no content field in prediction")
    }

    log.Printf("GeminiUsecase: Gemini API生レスポンス: %s", summaryText)

    geminiResponse, err := g.parseGeminiResponse(summaryText)
    if err != nil {
        log.Printf("GeminiUsecase: parseGeminiResponse エラー: %v", err)
        return &UserSummary{
            UserId:    userId,
            UserName:  posts[0].UserName,
            Summary:   summaryText,
            Interests: []string{"分析中..."},
            Personality: "分析中...",
        }, nil
    }
    log.Printf("GeminiUsecase: パースされたGeminiResponse: %+v", geminiResponse)

    return &UserSummary{
        UserId:    userId,
        UserName:  posts[0].UserName,
        Summary:   geminiResponse.Summary,
        Interests: geminiResponse.Interests,
        Personality: geminiResponse.Personality,
    }, nil
}

func (g *GeminiUsecase) parseGeminiResponse(text string) (*GeminiResponse, error) {
    jsonRegex := regexp.MustCompile(`\{[\s\S]*\}`)
    matches := jsonRegex.FindString(text)
    
    if matches == "" {
        log.Printf("parseGeminiResponse: レスポンスからJSONが見つかりません。テキスト: %s", text)
        return nil, fmt.Errorf("no JSON found in response")
    }
    
    var response GeminiResponse
    err := json.Unmarshal([]byte(matches), &response)
    if err != nil {
        log.Printf("parseGeminiResponse: JSONパースに失敗しました。マッチしたJSON: %s, エラー: %v", matches, err)
        return nil, fmt.Errorf("failed to parse JSON: %v", err)
    }
    
    return &response, nil
}