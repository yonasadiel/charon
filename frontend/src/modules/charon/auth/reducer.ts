import keyBy from 'lodash/keyBy';

import { PUT_USERS } from './action';
import { User } from './api';

export interface CharonAuthState {
  users: { [username: string]: User } | null,
};

const initialState: CharonAuthState = {
  users: null,
};

export function charonAuthReducer (state: CharonAuthState = initialState, action: any) {
  switch (action.type) {
    case PUT_USERS: {
      return {
        ...state,
        users: action.users === null ? null : keyBy(action.users, 'username'),
      };
    }
    default:
      return state;
  }
}
