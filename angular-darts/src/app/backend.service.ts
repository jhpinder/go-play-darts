import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ScoreboardResponse } from './scoreboard-response.Model';

@Injectable({
  providedIn: 'root'
})
export class BackendService {

  constructor(private http: HttpClient) { }

  throwDart(score: number) {
    const url = "localhost:8080/" + score;
    return this.http.get<ScoreboardResponse>(url);
  }
}
