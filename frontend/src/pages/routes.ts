export const ROUTE_HOME = '/';

export const ROUTE_LOGIN = ROUTE_HOME + 'login/';

export const ROUTE_VENUE = ROUTE_HOME + 'venue/';
export const ROUTE_VENUE_LIST = ROUTE_VENUE;

export const ROUTE_EXAM = ROUTE_HOME + 'exam/';
export const ROUTE_EVENT_LIST = ROUTE_EXAM;
export const ROUTE_EVENT_DETAIL = ROUTE_EXAM + ':eventId/';
export const ROUTE_EVENT_DETAIL_OVERVIEW = ROUTE_EVENT_DETAIL + 'overview/';
export const ROUTE_EVENT_QUESTION = ROUTE_EVENT_DETAIL + 'question/';
export const ROUTE_EVENT_QUESTION_DETAIL = ROUTE_EVENT_QUESTION + ':questionNumber/';
export const ROUTE_EVENT_QUESTION_EDIT = ROUTE_EVENT_DETAIL + 'question-editor/';
export const ROUTE_EVENT_QUESTION_EDIT_CREATE = ROUTE_EVENT_QUESTION_EDIT + 'new/';
