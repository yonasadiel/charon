import React from 'react';
import { matchPath } from 'react-router';
import { Link } from 'react-router-dom';


import { generateUrlWithParams } from '../../../../modules/util/routes';
import { ROUTE_EVENT_DETAIL_OVERVIEW, ROUTE_EVENT_QUESTION_DETAIL, ROUTE_EVENT_QUESTION_EDIT } from '../../../../pages/routes';
import './EventNavigation.scss';

interface EventNavigationProps {
  eventId: number | string;
  showEditQuestion: boolean;
  currentPath: string;
}

const EventNavigation = (props: EventNavigationProps) => {
  const { currentPath, eventId, showEditQuestion } = props;
  const urlParams = { eventId: eventId };

  const isActive = (path: string) =>
    !!matchPath(currentPath, { path, exact: true })
  return (
    <div className="event-navigation">
      <Link to={generateUrlWithParams(ROUTE_EVENT_DETAIL_OVERVIEW, urlParams)}>
        <div className={`nav-item ${isActive(ROUTE_EVENT_DETAIL_OVERVIEW) ? 'active' : ''}`}>
          <i className="fas fa-info-circle" /> Detail
        </div>
      </Link>
      <Link to={generateUrlWithParams(ROUTE_EVENT_QUESTION_DETAIL, { ...urlParams, questionId: 1 })}>
        <div className={`nav-item ${isActive(ROUTE_EVENT_QUESTION_DETAIL) ? 'active' : ''}`}>
          <i className="fas fa-book-open" /> Soal
        </div>
      </Link>
      {showEditQuestion && (
        <Link to={generateUrlWithParams(ROUTE_EVENT_QUESTION_EDIT, urlParams)}>
        <div className={`nav-item ${isActive(ROUTE_EVENT_QUESTION_EDIT) ? 'active' : ''}`}>
            <i className="fas fa-edit" /> Ubah Soal
          </div>
        </Link>
      )}
    </div>
  );
};

export default EventNavigation;
