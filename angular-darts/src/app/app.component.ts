import { Component } from '@angular/core';
import { BackendService } from './backend.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})

export class AppComponent {
  constructor(private backendService: BackendService) { }
  title = 'angular-darts';

  playerScores = {
    playerOne: 301,
    playerTwo: 301
  }

  throwDart(score: number) {
    this.backendService.throwDart(score).subscribe(response => {
      this.playerScores = response;
    });
  }
}