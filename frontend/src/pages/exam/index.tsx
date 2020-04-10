import React from 'react';
import { connect } from 'react-redux';
import { Route, Switch, Redirect } from 'react-router-dom';

import { User } from '../../modules/charon/auth/api';
import * as charonExamActions from '../../modules/charon/exam/action';
import * as sessionSelectors from '../../modules/session/selector';
import { AppState } from '../../modules/store';
import { ROUTE_LOGIN, ROUTE_EVENT_LIST, ROUTE_EVENT_DETAIL } from '../routes';
import EventListPage from './EventListPage';
import EventDetailPage from './_eventId/EventPage';

export interface ExamRouteProps {
  getEvents: () => void;
  user: User | null;
};

const ExamRoute = (props: ExamRouteProps) => {
  const { getEvents, user } = props;

  React.useEffect(() => { getEvents() }, [getEvents]);

  if (!user) {
    return <Redirect to={ROUTE_LOGIN} />;
  }

  return (
    <Switch>
      <Route path={ROUTE_EVENT_DETAIL}><EventDetailPage /></Route>
      <Route path={ROUTE_EVENT_LIST}><EventListPage user={user} /></Route>
    </Switch>
  );
};

const mapStateToProps = (state: AppState) => ({
  user: sessionSelectors.getUser(state),
});

const mapDispatchToProps = {
  getEvents: charonExamActions.getEvents,
};

export default connect(mapStateToProps, mapDispatchToProps)(ExamRoute);
