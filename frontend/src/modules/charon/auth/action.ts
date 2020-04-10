import { AxiosError, AxiosResponse } from 'axios';

import { putUser } from '../../session/action';
import { AppThunk } from '../../store';
import { CharonAPIError } from '../http';
import { User } from './api';

export function login(username: string, password: string): AppThunk<void> {
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

