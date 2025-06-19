package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"hackathon/dao"
	"hackathon/model"
	"google.golang.org/api/option"
	generativelanguage "google.golang.org/genproto/googleapis/ai/generativelanguage/v1beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GeminiUsecase struct {
	postDao dao.PostDao
	client  generativelanguage.GenerativeServiceClient
}

func NewGeminiUsecase(postDao dao.PostDao, apiKey string) (*GeminiUsecase, error) {
	ctx := context.Background()
	
	// Gemini APIクライアントの初期化
	conn, err := grpc.DialContext(ctx, "generativelanguage.googleapis.com:443", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.WithPerRPCCredentials(option.WithAPIKey(apiKey))))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Gemini API: %v", err)
	}

	client := generativelanguage.NewGenerativeServiceClient(conn)

	return &GeminiUsecase{
		postDao: postDao,
		client:  client,
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
	// ユーザーの投稿を取得
	posts, err := g.postDao.FindAllByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %v", err)
	}

	if len(posts) == 0 {
		return &UserSummary{
			UserId:    userId,
			UserName:  "Unknown",
			Summary:   "まだ投稿がありません",
			Interests: []string{},
			Personality: "投稿が少ないため、性格を分析できません",
		}, nil
	}

	// 投稿のテキストを結合
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

	// Gemini APIに送信するプロンプトを作成
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

	// Gemini APIを呼び出し
	request := &generativelanguage.GenerateContentRequest{
		Model: "models/gemini-1.5-flash",
		Contents: []*generativelanguage.Content{
			{
				Parts: []*generativelanguage.Part{
					{
						Text: prompt,
					},
				},
			},
		},
		GenerationConfig: &generativelanguage.GenerationConfig{
			Temperature:     0.7,
			TopP:           0.8,
			TopK:           40,
			MaxOutputTokens: 1000,
		},
	}

	response, err := g.client.GenerateContent(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini API")
	}

	// レスポンスを解析
	summaryText := response.Candidates[0].Content.Parts[0].Text
	
	// JSONを抽出してパース
	geminiResponse, err := g.parseGeminiResponse(summaryText)
	if err != nil {
		// パースに失敗した場合は、テキストをそのまま使用
		return &UserSummary{
			UserId:    userId,
			UserName:  posts[0].UserName,
			Summary:   summaryText,
			Interests: []string{"分析中..."},
			Personality: "分析中...",
		}, nil
	}
	
	return &UserSummary{
		UserId:    userId,
		UserName:  posts[0].UserName,
		Summary:   geminiResponse.Summary,
		Interests: geminiResponse.Interests,
		Personality: geminiResponse.Personality,
	}, nil
}

// Gemini APIのレスポンスからJSONを抽出してパースする
func (g *GeminiUsecase) parseGeminiResponse(text string) (*GeminiResponse, error) {
	// JSONブロックを抽出する正規表現
	jsonRegex := regexp.MustCompile(`\{[\s\S]*\}`)
	matches := jsonRegex.FindString(text)
	
	if matches == "" {
		return nil, fmt.Errorf("no JSON found in response")
	}
	
	var response GeminiResponse
	err := json.Unmarshal([]byte(matches), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}
	
	return &response, nil
} 