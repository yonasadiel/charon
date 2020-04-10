import { PUT_USER } from './action';
import { User } from '../charon/auth/api';

export interface SessionState {
  user: User | null;
};

const initialState: SessionState = {
  user: null,
};

export function sessionReducer (state: SessionState = initialState, action: any) {
  switch (action.type) {
    case PUT_USER: {
      return {
        ...state,
        user: action.user,
      };
    }
    default:
      return state;
  }
}