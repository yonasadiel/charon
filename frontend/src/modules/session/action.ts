import { User } from '../charon/auth/api';

export const PUT_USER = 'session/PUT_USER';
export const putUser = (user: User) => ({
  type: PUT_USER,
  user,
});
