import { Component } from '@angular/core';
import { DashboardComponent } from './dashboard/dashboard';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [DashboardComponent],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App {
  title = 'web-app';
}