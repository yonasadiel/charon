import { AxiosResponse } from 'axios';

import conf from '../../../conf';
import http from '../../http';

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
}

export default {
  login: (username: string, password: string) => http.post(`${conf.charonApiUrl}/auth/login/`, { username, password }),
  logout: () => http.post(`${conf.charonApiUrl}/auth/logout/`),

  getUsers: () => http.get(`${conf.charonApiUrl}/auth/user/`),
  createUser: (user: User) => http.post(`${conf.charonApiUrl}/auth/user/`, user),
}
