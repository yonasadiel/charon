import { AxiosError, AxiosResponse } from 'axios';

import { putUser } from '../../session/action';
import { AppThunk } from '../../store';
import { CharonAPIError, CharonFormError } from '../http';
import { User } from './api';

export const PUT_USERS = 'charon/auth/PUT_USERS';
export const putUsers = (users: User[] | null) => ({
  type: PUT_USERS,
  users,
});

export function login(username: string, password: string): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonAuthApi }) {
    return charonAuthApi.login(username, password)
      .then((res: AxiosResponse) => {
        const user: User = res.data;
        dispatch(putUser(user));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function logout(): AppThunk<void> {
  return async function (dispatch, _, { charonAuthApi }) {
    return charonAuthApi.logout()
      .then(() => {
        dispatch(putUser(null));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function getUsers(): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonAuthApi }) {
    dispatch(putUsers(null));
    return charonAuthApi.getUsers()
      .then((res: AxiosResponse) => {
        const users: User[] = res.data;
        dispatch(putUsers(users));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function createUser(user: User): AppThunk<Promise<void>> {
  return async function (_dispatch, _getState, { charonAuthApi }) {
    return charonAuthApi.createUser(user)
      .then(() => {
        return;
      })
      .catch((err: AxiosError) => {
        throw new CharonFormError(err);
      });
  };
};


