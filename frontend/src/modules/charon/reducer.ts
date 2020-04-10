import { combineReducers } from 'redux';

import { CharonExamState, charonExamReducer } from './exam/reducer';

export interface CharonState {
  exam: CharonExamState;
};

export const charonReducer = combineReducers<CharonState>({
  exam: charonExamReducer,
});
