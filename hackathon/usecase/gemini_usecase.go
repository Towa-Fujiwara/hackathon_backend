package usecase

import (
    "context"
    "encoding/json"
    "fmt"
    "regexp"
    "strings"
    "hackathon/dao"
    "cloud.google.com/go/vertexai/genai"
    "log" // logパッケージをインポート
)

type GeminiUsecase struct {
    postDao dao.PostDao
    client  *genai.Client
    model   *genai.GenerativeModel
}

func NewGeminiUsecase(postDao dao.PostDao, projectID, location, engineID string) (*GeminiUsecase, error) {
    ctx := context.Background()

    client, err := genai.NewClient(ctx, projectID, location)
    if err != nil {
        return nil, fmt.Errorf("failed to create Gemini client: %v", err)
    }

    model := client.GenerativeModel(engineID)

    return &GeminiUsecase{
        postDao: postDao,
        client:  client,
        model:   model,
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
    log.Printf("GeminiUsecase: GenerateUserSummaryで受信したuserId: %s", userId) // GenerateUserSummaryで受け取ったuserIdをログ出力

    posts, err := g.postDao.FindAllByUserId(userId)
    if err != nil {
        log.Printf("GeminiUsecase: g.postDao.FindAllByUserId(%s) エラー: %v", userId, err) // g.postDao.FindAllByUserIdのエラーをログ出力
        return nil, fmt.Errorf("failed to get user posts: %v", err)
    }
    log.Printf("GeminiUsecase: g.postDao.FindAllByUserId(%s) で取得した投稿数: %d", userId, len(posts)) // 取得した投稿数をログ出力

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

    resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        log.Printf("GeminiUsecase: Gemini API GenerateContent エラー: %v", err) // Gemini API GenerateContentのエラーをログ出力
        return nil, fmt.Errorf("failed to generate content: %v", err)
    }

    if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
        log.Print("GeminiUsecase: Gemini APIからのレスポンス候補がない、またはコンテンツが空です") // Gemini APIからのレスポンス候補がない、またはコンテンツが空の場合をログ出力
        return nil, fmt.Errorf("no response from Gemini API")
    }

    var summaryText string
    for _, part := range resp.Candidates[0].Content.Parts {
        if txt, ok := part.(genai.Text); ok {
            summaryText += string(txt)
        }
    }
    log.Printf("GeminiUsecase: Gemini API生レスポンス: %s", summaryText) // Gemini APIからの生レスポンステキストをログ出力

    geminiResponse, err := g.parseGeminiResponse(summaryText)
    if err != nil {
        log.Printf("GeminiUsecase: parseGeminiResponse エラー: %v", err) // parseGeminiResponseのエラーをログ出力
        return &UserSummary{
            UserId:    userId,
            UserName:  posts[0].UserName,
            Summary:   summaryText,
            Interests: []string{"分析中..."},
            Personality: "分析中...",
        }, nil
    }
    log.Printf("GeminiUsecase: パースされたGeminiResponse: %+v", geminiResponse) // パースされたGeminiResponseをログ出力

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
        log.Printf("parseGeminiResponse: レスポンスからJSONが見つかりません。テキスト: %s", text) // レスポンスからJSONが見つからない場合をログ出力
        return nil, fmt.Errorf("no JSON found in response")
    }
    
    var response GeminiResponse
    err := json.Unmarshal([]byte(matches), &response)
    if err != nil {
        log.Printf("parseGeminiResponse: JSONパースに失敗しました。マッチしたJSON: %s, エラー: %v", matches, err) // JSONパースに失敗した場合をログ出力
        return nil, fmt.Errorf("failed to parse JSON: %v", err)
    }
    
    return &response, nil
}