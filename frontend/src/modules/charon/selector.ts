import { AppState } from '../store';
import { CharonState } from './reducer';

export function getCharonState(state: AppState): CharonState {
  return state.charon;
}
