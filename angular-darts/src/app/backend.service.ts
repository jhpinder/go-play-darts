import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { GameStatusResponse } from './game-status-response.Model';
import { ThrowDartRequest } from './throw-dart-request.Model'

@Injectable({
  providedIn: 'root'
})
export class BackendService {

  constructor(private http: HttpClient) { }

  throwDart(score: number, gameID: number, playerName: string) {
    const url = "http://192.168.50.244:9090/game/" + gameID;
    const body: ThrowDartRequest = {
      Score: score,
      PlayerName: playerName
    }
    return this.http.post<GameStatusResponse>(url, body);
  }

  newGame() {
    const url = "http://192.168.50.244:9090/game/new";
    return this.http.post<GameStatusResponse>(url, null);
  }

  currentGameState(gameID: number) {
    const url = "http://192.168.50.244:9090/game/" + gameID;
    return this.http.get<GameStatusResponse>(url);
  }

  getLatestGame() {
    const url = "http://192.168.50.244:9090/game/latest";
    return this.http.get<GameStatusResponse>(url);
  }
}
