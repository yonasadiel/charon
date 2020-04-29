import { PUT_USER, PUT_PARTICIPATION_KEY } from './action';
import { User } from '../charon/auth/api';

export interface SessionState {
  user: User | null;
  participationKey: { [eventSlug: string]: string };
};

const initialState: SessionState = {
  user: null,
  participationKey: {},
};

export function sessionReducer (state: SessionState = initialState, action: any) {
  switch (action.type) {
    case PUT_USER: {
      return {
        ...state,
        user: action.user,
      };
    }
    case PUT_PARTICIPATION_KEY: {
      const newParticipationKey = Object.assign({}, state.participationKey);
      newParticipationKey[action.eventSlug] = action.participationKey;
      return {
        ...state,
        participationKey: newParticipationKey,
      };
    }
    default:
      return state;
  }
}
