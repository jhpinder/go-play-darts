import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { GameStatus } from './game-status.Model';
import { ThrowDartRequest } from './throw-dart-request.Model'

@Injectable({
  providedIn: 'root'
})
export class BackendService {

  constructor(private http: HttpClient) { }

  throwDart(score: number, gameID: string) {
    const url = "http://192.168.50.244:9090/game/" + gameID;
    const body = score;
    return this.http.post<GameStatus>(url, body);
  }

  newGame(playerList: string[]) {
    const url = "http://192.168.50.244:9090/game/new";
    return this.http.post<GameStatus>(url, playerList);
  }

  currentGameState(gameID: number) {
    const url = "http://192.168.50.244:9090/game/" + gameID;
    return this.http.get<GameStatus>(url);
  }

  getLatestGame() {
    const url = "http://192.168.50.244:9090/game/latest";
    return this.http.get<GameStatus>(url);
  }
}
