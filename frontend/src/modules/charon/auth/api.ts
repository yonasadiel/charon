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
  role: 'local' | 'participant' | 'organizer' | 'admin';
};

export interface CharonAuthApi {
  login: (username: string, password: string) => Promise<AxiosResponse<any>>;
}

export default {
  login: (username: string, password: string) => {
    return http.post(`${conf.charonApiUrl}/auth/login/`, { username, password });
  },
}
