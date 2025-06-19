# Gemini機能の使用方法

このプロジェクトでは、ユーザーの投稿からAIが自動的にサマリーを生成するGemini機能を搭載しています。

## 機能概要

- ユーザーの投稿内容を分析
- 投稿内容の要約を生成
- ユーザーの興味・関心を抽出
- 投稿から読み取れる性格・特徴を分析

## セットアップ

### 1. Gemini APIキーの取得

1. [Google AI Studio](https://makersuite.google.com/app/apikey)にアクセス
2. APIキーを作成
3. 作成したAPIキーをコピー

### 2. 環境変数の設定

```bash
export GEMINI_API_KEY="your_gemini_api_key_here"
```

または、`.env`ファイルに追加：

```
GEMINI_API_KEY=your_gemini_api_key_here
```

### 3. 依存関係のインストール

```bash
go mod tidy
```

## APIエンドポイント

### 1. 自分の投稿からサマリーを生成

**エンドポイント:** `GET /api/users/me/summary`

**認証:** 必要（Firebase認証）

**レスポンス例:**
```json
{
  "userId": "user123",
  "userName": "田中太郎",
  "summary": "技術系の投稿が多く、特にGo言語とWeb開発について詳しく発信しています。新しい技術の学習意欲が高く、コミュニティ活動にも積極的に参加している様子が伺えます。",
  "interests": ["プログラミング", "Go言語", "Web開発", "技術学習"],
  "personality": "技術への関心が高く、学習意欲旺盛で、知識を共有することに積極的な性格"
}
```

### 2. 特定ユーザーの投稿からサマリーを生成

**エンドポイント:** `GET /api/users/{userId}/summary`

**認証:** 不要

**パラメータ:**
- `userId`: 分析対象のユーザーID

**レスポンス例:**
```json
{
  "userId": "user456",
  "userName": "佐藤花子",
  "summary": "料理と旅行に関する投稿が中心で、特に海外旅行での体験談や現地の料理について詳しく発信しています。",
  "interests": ["料理", "旅行", "海外文化", "写真"],
  "personality": "好奇心旺盛で、新しい体験を楽しむことが好きな性格"
}
```

## エラーハンドリング

- APIキーが設定されていない場合、機能は無効化されます
- 投稿が存在しない場合、適切なメッセージが返されます
- Gemini APIの呼び出しに失敗した場合、エラーメッセージが返されます

## 注意事項

- Gemini APIの利用には料金が発生する場合があります
- 大量のリクエストを送信する場合は、レート制限に注意してください
- 投稿内容の分析結果は、AIによる推測であり、100%正確ではない場合があります

## トラブルシューティング

### よくある問題

1. **APIキーエラー**
   - 環境変数`GEMINI_API_KEY`が正しく設定されているか確認
   - APIキーが有効かどうか確認

2. **投稿が見つからない**
   - 指定したユーザーIDが正しいか確認
   - ユーザーに投稿が存在するか確認

3. **分析結果が不正確**
   - 投稿内容が十分にあるか確認
   - テキスト投稿が含まれているか確認

## 開発者向け情報

### アーキテクチャ

- **Controller**: `controller/gemini_controller.go`
- **Usecase**: `usecase/gemini_usecase.go`
- **Model**: `model/post.go`（既存のPostモデルを使用）

### カスタマイズ

Gemini APIのプロンプトや分析内容をカスタマイズする場合は、`usecase/gemini_usecase.go`の`GenerateUserSummary`メソッドを編集してください。 