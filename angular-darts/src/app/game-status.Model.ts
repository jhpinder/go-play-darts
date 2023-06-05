export interface GameStatus {
  Scores: Map<string, number>;
  GameID: string
  OrderedPlayers: string[]
  CurrentPlayerIndex: number
}
