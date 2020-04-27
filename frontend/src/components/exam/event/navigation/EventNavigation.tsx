import React from 'react';
import { matchPath } from 'react-router';
import { Link } from 'react-router-dom';


import { generateUrlWithParams } from '../../../../modules/util/routes';
import {
  ROUTE_EVENT_OVERVIEW,
  ROUTE_EVENT_PARTICIPATION,
  ROUTE_EVENT_QUESTION_DETAIL,
  ROUTE_EVENT_QUESTION_EDIT,
} from '../../../../pages/routes';
import './EventNavigation.scss';

interface EventNavigationProps {
  currentPath: string;
  eventSlug: string;
  hasEditPermission: boolean;
}

const EventNavigation = (props: EventNavigationProps) => {
  const { currentPath, eventSlug, hasEditPermission } = props;
  const urlParams = { eventSlug };

  const isActive = (path: string) =>
    !!matchPath(currentPath, { path, exact: true })
  return (
    <div className="event-navigation">
      <Link to={generateUrlWithParams(ROUTE_EVENT_OVERVIEW, urlParams)}>
        <div className={`nav-item ${isActive(ROUTE_EVENT_OVERVIEW) ? 'active' : ''}`}>
          <i className="fas fa-info-circle" /> Detail
        </div>
      </Link>
      {hasEditPermission && (
        <Link to={generateUrlWithParams(ROUTE_EVENT_PARTICIPATION, urlParams)}>
        <div className={`nav-item ${isActive(ROUTE_EVENT_PARTICIPATION) ? 'active' : ''}`}>
            <i className="fas fa-user-friends" /> Peserta
          </div>
        </Link>
      )}
      <Link to={generateUrlWithParams(ROUTE_EVENT_QUESTION_DETAIL, { ...urlParams, questionId: 1 })}>
        <div className={`nav-item ${isActive(ROUTE_EVENT_QUESTION_DETAIL) ? 'active' : ''}`}>
          <i className="fas fa-book-open" /> Soal
        </div>
      </Link>
      {hasEditPermission && (
        <Link to={generateUrlWithParams(ROUTE_EVENT_QUESTION_EDIT, urlParams)}>
        <div className={`nav-item ${isActive(ROUTE_EVENT_QUESTION_EDIT) ? 'active' : ''}`}>
            <i className="fas fa-edit" /> Ubah Soal
          </div>
        </Link>
      )}
    </div>
  );
};

EventNavigation.defaultProps = {
  hasEditPermission: false,
}

export default EventNavigation;
