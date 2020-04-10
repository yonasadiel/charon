import keyBy from 'lodash/keyBy';

import { PUT_EVENTS, PUT_QUESTIONS } from './action';
import { Event } from './api';

export interface CharonExamState {
  events: { [id: number]: Event } | null,
};

const initialState: CharonExamState = {
  events: null,
};

export function charonExamReducer (state: CharonExamState = initialState, action: any) {
  switch (action.type) {
    case PUT_EVENTS: {
      return {
        ...state,
        events: action.events === null ? null : keyBy(action.events, 'id'),
      };
    }
    case PUT_QUESTIONS: {
      const oldEvent = !!state.events ? state.events[action.eventId] : null;
      const newEvent = Object.assign({}, !!oldEvent ? oldEvent : {} as Event);
      newEvent.questions = action.questions === null ? null : keyBy(action.questions, 'id');
      return {
        ...state,
        events: {
          ...state.events,
          [action.eventId]: newEvent,
        },
      };
    }
    default:
      return state;
  }
}