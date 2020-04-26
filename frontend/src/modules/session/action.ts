import { User } from '../charon/auth/api';

export const PUT_USER = 'session/PUT_USER';
export const putUser = (user: User | null) => ({
  type: PUT_USER,
  user,
});
