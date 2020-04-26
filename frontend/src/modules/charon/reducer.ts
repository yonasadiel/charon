import { combineReducers } from 'redux';

import { CharonAuthState, charonAuthReducer } from './auth/reducer';
import { CharonExamState, charonExamReducer } from './exam/reducer';

export interface CharonState {
  auth: CharonAuthState;
  exam: CharonExamState;
};

export const charonReducer = combineReducers<CharonState>({
  auth: charonAuthReducer,
  exam: charonExamReducer,
});
