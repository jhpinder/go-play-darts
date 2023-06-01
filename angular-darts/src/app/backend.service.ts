import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { GameStatusResponse } from './game-status-response.Model';

@Injectable({
  providedIn: 'root'
})
export class BackendService {

  constructor(private http: HttpClient) { }

  throwDart(score: number) {
    const url = "http://192.168.50.244:9090/" + score;
    return this.http.get<GameStatusResponse>(url);
  }

  newGame() {
    const url = "http://192.168.50.244:9090/restart";
    return this.http.get<GameStatusResponse>(url);
  }

  currentGameState() {
    const url = "http://192.168.50.244:9090";
    return this.http.get<GameStatusResponse>(url);
  }
}
