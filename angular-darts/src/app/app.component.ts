import { Component } from '@angular/core';
import { map } from 'rxjs';
import { BackendService } from './backend.service';
import { GameStatus as GameStatus } from './game-status.Model';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})

export class AppComponent {
  constructor(private backendService: BackendService) { }
  title = 'angular-darts';

  showNewGame = true;
  gameStatus!: GameStatus;
  playerList: string[] = [];
  currentPlayerID: string = "";

  ngOnInit() {

  }

  throwDart(score: number) {
    this.backendService.throwDart(score, this.gameStatus.GameID)
      .pipe(map((response: GameStatus) => {
        let toReturn: GameStatus = {
          Scores: new Map<string, number>(Object.entries(response.Scores)),
          GameID: response.GameID,
          OrderedPlayers: response.OrderedPlayers,
          CurrentPlayerIndex: response.CurrentPlayerIndex
        }
        return toReturn;
      })).subscribe(response => {
        this.gameStatus = response;
      });
  }

  newGame() {
    this.backendService.newGame(this.playerList)
      .pipe(map((response: GameStatus) => {
        let toReturn: GameStatus = {
          Scores: new Map<string, number>(Object.entries(response.Scores)),
          GameID: response.GameID,
          OrderedPlayers: response.OrderedPlayers,
          CurrentPlayerIndex: response.CurrentPlayerIndex
        }
        return toReturn;
      })).subscribe(response => {
        this.gameStatus = response;
        this.showNewGame = false;
      });
  }

  addPlayer() {
    if (this.currentPlayerID.length < 1
      || this.playerList.find(element => element == this.currentPlayerID)) {
      this.currentPlayerID = "";
      return;
    }
    this.playerList.push(this.currentPlayerID);
    this.currentPlayerID = "";
  }

  endGame() {
    this.showNewGame = true;
  }
}
