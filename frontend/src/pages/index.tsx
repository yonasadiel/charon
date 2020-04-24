import * as React from 'react';
import { Route, Switch } from 'react-router-dom';

import Navbar from '../components/navigation/Navbar';
import ExamRoute from './exam';
import VenueListPage from './venue/VenueListPage';
import LoginPage from './login/LoginPage';
import HomePage from './HomePage';
import { ROUTE_EXAM, ROUTE_HOME, ROUTE_LOGIN, ROUTE_VENUE } from './routes';

const RouteComponent = () => {
  return (
    <>
      <Navbar />
      <Switch>
        <Route path={ROUTE_LOGIN} component={LoginPage} />
        <Route path={ROUTE_EXAM} component={ExamRoute} />
        <Route path={ROUTE_VENUE} component={VenueListPage} />
        <Route path={ROUTE_HOME} component={HomePage} />
      </Switch>
    </>
  );
}
export default RouteComponent;
