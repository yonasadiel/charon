import keyBy from 'lodash/keyBy';

import { PUT_EVENTS, PUT_QUESTIONS, PUT_VENUES } from './action';
import { Event, Venue } from './api';

export interface CharonExamState {
  venues: { [id: number]: Venue } | null,
  events: { [id: number]: Event } | null,
};

const initialState: CharonExamState = {
  venues: null,
  events: null,
};

export function charonExamReducer (state: CharonExamState = initialState, action: any) {
  switch (action.type) {
    case PUT_VENUES: {
      return {
        ...state,
        venues: action.venues === null ? null : keyBy(action.venues, 'id'),
      };
    }
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
