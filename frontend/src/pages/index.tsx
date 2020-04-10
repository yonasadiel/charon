import * as React from 'react';
import { Route, Switch } from 'react-router-dom';

import Navbar from '../components/navigation/Navbar';
import ExamRoute from './exam';
import LoginPage from './login/LoginPage';
import HomePage from './HomePage';
import { ROUTE_LOGIN, ROUTE_EXAM, ROUTE_HOME } from './routes';

const RouteComponent = () => {
  return (
    <>
      <Navbar />
      <Switch>
        <Route path={ROUTE_LOGIN} component={LoginPage} />
        <Route path={ROUTE_EXAM} component={ExamRoute} />
        <Route path={ROUTE_HOME} component={HomePage} />
      </Switch>
    </>
  );
}
export default RouteComponent;
