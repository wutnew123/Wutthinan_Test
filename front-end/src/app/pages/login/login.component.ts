import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class Login {
  username = '';
  password = '';
  errorMessage = '';
  loading = false;

  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  onLogin(): void {
    this.errorMessage = '';
    if (!this.username.trim() || !this.password) {
      this.errorMessage = 'กรุณากรอกข้อมูลให้ครบทุกช่อง';
      return;
    }

    if (!this.username.trim() || !this.password.trim()) {
      this.errorMessage = 'กรุณากรอก Username และ Password';
      return;
    }

    this.loading = true;

    this.authService.login({
      username: this.username,
      password: this.password
    }).subscribe({
      next: (response) => {
        this.loading = false;

        if (response.status === 'success') {
          this.router.navigate(['/home']);
        } else {
          this.errorMessage = response.message || 'เข้าสู่ระบบไม่สำเร็จ';
        }
      },
      error: () => {
        this.loading = false;
        this.errorMessage =
          'เข้าสู่ระบบไม่สำเร็จ กรุณาตรวจสอบข้อมูล';
      }
    });
  }
}
