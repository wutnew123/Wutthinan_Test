import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class Register {
  name = '';
  password = '';
  confirmPassword = '';

  errorMessage = '';
  successMessage = '';
  loading = false;

  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  onRegister(): void {
    this.errorMessage = '';
    this.successMessage = '';

    if (
      !this.name.trim() ||
      !this.password ||
      !this.confirmPassword
    ) {
      this.errorMessage = 'กรุณากรอกข้อมูลให้ครบทุกช่อง';
      return;
    }

    if (this.name.trim().length < 3) {
      this.errorMessage = 'Username ต้องมีอย่างน้อย 3 ตัวอักษร';
      return;
    }

    if (this.password.length < 8) {
      this.errorMessage = 'Password ต้องมีอย่างน้อย 8 ตัวอักษร';
      return;
    }

    if (this.password !== this.confirmPassword) {
      this.errorMessage = 'Password และ Confirm Password ไม่ตรงกัน';
      return;
    }

    this.loading = true;

    this.authService.register({
      name: this.name.trim(),
      password: this.password
    }).subscribe({
      next: () => {
        this.loading = false;
        this.successMessage = 'สมัครสมาชิกสำเร็จ';

        this.router.navigate(['/login'], {
          state: { message: 'สมัครสมาชิกสำเร็จ กรุณาเข้าสู่ระบบ' }
        });
      },
      error: (error) => {
        this.loading = false;

        this.errorMessage =
          error.error?.message ??
          'สมัครสมาชิกไม่สำเร็จ กรุณาตรวจสอบข้อมูล';
      }
    });
  }
}
