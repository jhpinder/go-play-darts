import { Component } from '@angular/core';
import { BackendService } from './backend.service';
import { GameStatusResponse } from './game-status-response.Model';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})

export class AppComponent {
  constructor(private backendService: BackendService) { }
  title = 'angular-darts';

  gameStatus!: GameStatusResponse;

  ngOnInit() {
    this.backendService.getLatestGame().subscribe(response => {
      this.gameStatus = response;
    });
  }

  throwDart(score: number) {
    this.backendService.throwDart(score, this.gameStatus.GameID, "Bob").subscribe(response => {
      this.gameStatus = response;
    });
  }

  newGame() {
    this.backendService.newGame().subscribe(response => {
      this.gameStatus = response;
    });
  }
}
