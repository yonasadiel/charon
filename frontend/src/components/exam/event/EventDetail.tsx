import React from 'react';
import { Card } from 'react-hephaestus';

import { Event } from '../../../modules/charon/exam/api';
import './EventDetail.scss';

interface EventDetailProps {
  event: Event;
};

const EventDetail = (props: EventDetailProps) => {
  const { event } = props;
  const diff = event.endsAt.getTime() - event.startsAt.getTime();
  const durationDays = Math.floor(diff / (24 * 60 * 60 * 1000));
  const durationHours = Math.floor(diff / (60 * 60 * 1000)) % 24;
  const durationMinutes = Math.floor(diff / (60 * 1000)) % 60;
  const durationSeconds = Math.floor(diff / (1000)) % 60;
  let durationText = '';
  if (durationDays > 0) {
    durationText = `${durationDays} hari ${durationHours} jam ${durationMinutes} menit ${durationSeconds} detik`;
  } else if (durationHours > 0) {
    durationText = `${durationHours} jam ${durationMinutes} menit ${durationSeconds} detik`;
  } else if (durationMinutes > 0) {
    durationText = `${durationMinutes} menit ${durationSeconds} detik`;
  } else if (durationSeconds > 0) {
    durationText = `${durationSeconds} detik`;
  }
  return (
    <Card className="event-detail">
      <p>{event.description}</p>
      <ul>
        <li>Waktu mulai: {event.startsAt.toLocaleDateString()} {event.startsAt.toLocaleTimeString()}</li>
        <li>Durasi: {durationText}</li>
      </ul>
    </Card>
  );
};

export default EventDetail;