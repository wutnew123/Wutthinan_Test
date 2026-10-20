import { Injectable, PLATFORM_ID, inject } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, finalize, tap } from 'rxjs';
import { environment } from '../../environments/environment';

export interface RegisterRequest {
  name: string;
  password: string;
}
export interface LoginRequest {
  username: string;
  password: string;
}
 
export interface LoginResponse {
  status: string;
  message: string;
  data: {
    access_token: string;
    datas: {
      name: string;
      username: string;
    };
    refresh_token: string;
  };
}

export interface LogoutResponse {
  status: string;
  message: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly apiUrl = environment.apiUrl;
  private readonly tokenKey = 'access_token';
  private readonly platformId = inject(PLATFORM_ID);

  constructor(private http: HttpClient) {}
  
  register(data: RegisterRequest) {
    return this.http.post(
      `${this.apiUrl}/auth/register`,
      data
    );
  }
  
  login(data: LoginRequest): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>(`${this.apiUrl}/auth/login`, data)
      .pipe(
        tap(response => {
 
          if (response.status === 'success') {
            const storage = this.getStorage();

            storage?.setItem(this.tokenKey, response.data.access_token);
            storage?.setItem('user_name', response.data.datas.name);
            storage?.setItem('refresh_token', response.data.refresh_token);
          }
        })
      );
  }

  logout(): Observable<LogoutResponse> {
    const storage = this.getStorage();
    const accessToken = storage?.getItem(this.tokenKey);
    const refreshToken = storage?.getItem('refresh_token');
    let headers = new HttpHeaders();

    if (accessToken) {
      headers = headers.set('Authorization', `Bearer ${accessToken}`);
    }

    return this.http.post<LogoutResponse>(
      `${this.apiUrl}/auth/logout`,
      { refresh_token: refreshToken },
      { headers }
    ).pipe(
      finalize(() => this.clearSession())
    );
  }

  private clearSession(): void {
    const storage = this.getStorage();

    storage?.removeItem(this.tokenKey);
    storage?.removeItem('refresh_token');
    storage?.removeItem('user_name');
  }

  isLoggedIn(): boolean {
    return !!this.getStorage()?.getItem(this.tokenKey);
  }

  getToken(): string | null {
    return this.getStorage()?.getItem(this.tokenKey) ?? null;
  }

  getUserName(): string | null {
    return this.getStorage()?.getItem('user_name') ?? null;
  }

  hasRefreshToken(): boolean {
    return !!this.getStorage()?.getItem('refresh_token');
  }

  private getStorage(): Storage | null {
    return isPlatformBrowser(this.platformId) ? sessionStorage : null;
  }
  getProfile() {
  return this.http.get<any>(`${this.apiUrl}/auth/me`);
 }
}
