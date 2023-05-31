import { Component } from '@angular/core';
import { BackendService } from './backend.service';
import { ScoreboardResponse } from './scoreboard-response.Model';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})

export class AppComponent {
  constructor(private backendService: BackendService) { }
  title = 'angular-darts';

  playerScores!: ScoreboardResponse;
  backendError!: boolean;

  ngOnInit() {
    this.backendService.currentGameState().subscribe(response => {
      this.playerScores = response;
      this.backendError = false;
    }, error => {
      this.backendError = true
    });
  }

  throwDart(score: number) {
    this.backendService.throwDart(score).subscribe(response => {
      this.playerScores = response;
    });
  }

  newGame() {
    this.backendService.newGame().subscribe(response => {
      this.playerScores = response;
    });
  }
}
