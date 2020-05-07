import React from 'react';
import { matchPath } from 'react-router';

import { generateUrlWithParams } from '../../../../modules/util/routes';
import {
  ROUTE_EVENT_OVERVIEW,
  ROUTE_EVENT_PARTICIPATION,
  ROUTE_EVENT_PARTICIPATION_STATUS,
  ROUTE_EVENT_QUESTION_DETAIL,
  ROUTE_EVENT_QUESTION_EDIT,
  ROUTE_EVENT_SYNC,
  ROUTE_EVENT_DECRYPT,
} from '../../../../pages/routes';
import EventNavigationItem from './EventNavigationItem';
import './EventNavigation.scss';

interface EventNavigationProps {
  currentPath: string;
  eventSlug: string;
  menus: string[];
}

const EventNavigation = (props: EventNavigationProps) => {
  const { currentPath, eventSlug, menus } = props;
  const urlParams = { eventSlug };

  const isActive = (path: string) =>
    !!matchPath(currentPath, { path, exact: true })
  return (
    <div className="event-navigation">
      {(menus.includes(ROUTE_EVENT_OVERVIEW)) && (
        <EventNavigationItem
          icon="fa-info-circle"
          isActive={isActive(ROUTE_EVENT_OVERVIEW)}
          route={generateUrlWithParams(ROUTE_EVENT_OVERVIEW, urlParams)}
          text="Detail" /> )}
      {(menus.includes(ROUTE_EVENT_PARTICIPATION)) && (
        <EventNavigationItem
          icon="fa-user-friends"
          isActive={isActive(ROUTE_EVENT_PARTICIPATION)}
          route={generateUrlWithParams(ROUTE_EVENT_PARTICIPATION, urlParams)}
          text="Peserta" /> )}
      {(menus.includes(ROUTE_EVENT_PARTICIPATION_STATUS)) && (
        <EventNavigationItem
          icon="fa-user-lock"
          isActive={isActive(ROUTE_EVENT_PARTICIPATION_STATUS)}
          route={generateUrlWithParams(ROUTE_EVENT_PARTICIPATION_STATUS, urlParams)}
          text="Status Peserta" /> )}
      {(menus.includes(ROUTE_EVENT_QUESTION_DETAIL)) && (
        <EventNavigationItem
          icon="fa-book-open"
          isActive={isActive(ROUTE_EVENT_QUESTION_DETAIL)}
          route={generateUrlWithParams(ROUTE_EVENT_QUESTION_DETAIL, { ...urlParams, questionNumber: 1 })}
          text="Soal" /> )}
      {(menus.includes(ROUTE_EVENT_QUESTION_EDIT)) && (
        <EventNavigationItem
          icon="fa-edit"
          isActive={isActive(ROUTE_EVENT_QUESTION_EDIT)}
          route={generateUrlWithParams(ROUTE_EVENT_QUESTION_EDIT, urlParams)}
          text="Ubah Soal" /> )}
      {(menus.includes(ROUTE_EVENT_SYNC)) && (
        <EventNavigationItem
          icon="fa-sync"
          isActive={isActive(ROUTE_EVENT_SYNC)}
          route={generateUrlWithParams(ROUTE_EVENT_SYNC, urlParams)}
          text="Sinkronisasi" /> )}
      {(menus.includes(ROUTE_EVENT_DECRYPT)) && (
        <EventNavigationItem
          icon="fa-key"
          isActive={isActive(ROUTE_EVENT_DECRYPT)}
          route={generateUrlWithParams(ROUTE_EVENT_DECRYPT, urlParams)}
          text="Dekripsi" /> )}
    </div>
  );
};

EventNavigation.defaultProps = {
  hasEditPermission: false,
}

export default EventNavigation;
