import parseInt from 'lodash/parseInt';
import React from 'react';
import { connect } from 'react-redux';
import { withRouter, RouteComponentProps, Switch, Route, Redirect } from 'react-router-dom';

import { User, USER_ROLE } from '../../../modules/charon/auth/api';
import EventDetail from '../../../components/exam/event/EventDetail';
import EventNavigation from '../../../components/exam/event/navigation/EventNavigation';
import { Event, Question } from '../../../modules/charon/exam/api';
import * as charonExamActions from '../../../modules/charon/exam/action';
import * as charonExamSelectors from '../../../modules/charon/exam/selector';
import * as charonSessionSelectors from '../../../modules/session/selector';
import { generateUrlWithParams } from '../../../modules/util/routes';
import { AppState } from '../../../modules/store';
import {
  ROUTE_EVENT,
  ROUTE_EVENT_OVERVIEW,
  ROUTE_EVENT_QUESTION,
  ROUTE_EVENT_QUESTION_DETAIL,
  ROUTE_EVENT_QUESTION_EDIT,
  ROUTE_EVENT_QUESTION_EDIT_CREATE,
} from '../../routes';
import EventQuestionCreatePage from './question-editor/create/QuestionCreatePage';
import QuestionEditorPage from './question-editor/QuestionEditorPage';
import QuestionPage from './question/QuestionPage';
import './EventPage.scss';

interface EventDetailPageProps extends RouteComponentProps<{ eventId: string }> {
  event: Event | null;
  user: User | null;
  createQuestion: (eventId: number, question: Question) => Promise<void>;
  getQuestionsOfEvent: (eventId: number) => Promise<void>;
};

const canEditQuestion = (user: User) => (user.role === USER_ROLE.ADMIN || user.role === USER_ROLE.ORGANIZER);

const EventDetailPage = (props: EventDetailPageProps) => {
  const {
    location: { pathname },
    match: { params: { eventId } },
    event,
    user,
  } = props;

  React.useEffect(() => { document.title = event?.title || 'Ujian'; }, [event]);

  if (!event) return <p>Loading</p>;

  const showEditQuestion = !!user ? canEditQuestion(user) : false;
  const eventDetailLink = generateUrlWithParams(ROUTE_EVENT_OVERVIEW, { eventId: event.id });
  const questionDetailLink = generateUrlWithParams(ROUTE_EVENT_QUESTION_DETAIL, { eventId: event.id, questionNumber: 1 });

  return (
    <div className="event-detail-page">
      <div className="titlebar">
        <h1>{event.title}</h1>
      </div>
      <div className="menubar">
        <EventNavigation currentPath={pathname} eventId={eventId} showEditQuestion={showEditQuestion} />
      </div>
      <div className="content">
        <Switch>
          <Route path={ROUTE_EVENT_QUESTION_EDIT_CREATE}>
            {showEditQuestion
              ? <EventQuestionCreatePage event={event} />
              : <Redirect to={questionDetailLink} />}
          </Route>
          <Route path={ROUTE_EVENT_QUESTION_EDIT}>
            {showEditQuestion
              ? <QuestionEditorPage eventId={event.id} />
              : <Redirect to={questionDetailLink} />}
          </Route>
          <Route path={ROUTE_EVENT_QUESTION_DETAIL}>
            <QuestionPage eventId={event.id} />
          </Route>
          <Route path={ROUTE_EVENT_QUESTION}>
            <Redirect to={questionDetailLink} />
          </Route>
          <Route path={ROUTE_EVENT_OVERVIEW}>
            <EventDetail event={event} />
          </Route>
          <Route path={ROUTE_EVENT}>
            <Redirect to={eventDetailLink} />
          </Route>
        </Switch>
      </div>
    </div>
  );
};

const mapStateToProps = (state: AppState, props: RouteComponentProps<{ eventId: string }>) => ({
  event: charonExamSelectors.getEvent(state, parseInt(props.match.params.eventId)),
  user: charonSessionSelectors.getUser(state),
});

const mapDispatchToProps = {
  createQuestion: charonExamActions.createQuestion,
  getQuestionsOfEvent: charonExamActions.getQuestionsOfEvent,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(EventDetailPage));
