import { AxiosResponse } from 'axios';

import http from '../http';

export const USER_ROLE = {
  LOCAL: 'local',
  PARTICIPANT: 'participant',
  ORGANIZER: 'organizer',
  ADMIN: 'admin',
};

export interface User {
  name: string;
  username: string;
  password: string;
  role: 'local' | 'participant' | 'organizer' | 'admin';
};

export interface CharonAuthApi {
  login: (username: string, password: string) => Promise<AxiosResponse<void>>;
  logout: () => Promise<AxiosResponse<void>>;

  getUsers: () => Promise<AxiosResponse<User[]>>;
  createUser: (user: User) => Promise<AxiosResponse<void>>;
  lockUser: (username: string) => Promise<AxiosResponse<void>>;
  unlockUser: (username: string) => Promise<AxiosResponse<void>>;
}

export default {
  login: (username: string, password: string) => http.post(`/auth/login/`, { username, password }),
  logout: () => http.post(`/auth/logout/`),

  getUsers: () => http.get(`/auth/user/`),
  createUser: (user: User) => http.post(`/auth/user/`, user),
  lockUser: (username: string) => http.post(`/auth/user/${username}/lock/`),
  unlockUser: (username: string) => http.post(`/auth/user/${username}/unlock/`),
}
