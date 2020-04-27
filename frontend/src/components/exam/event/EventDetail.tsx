import React from 'react';
import { Card } from 'react-hephaestus';

import { Event } from '../../../modules/charon/exam/api';
import { durationText } from '../../../modules/util/time';
import './EventDetail.scss';

interface EventDetailProps {
  event: Event | null;
};

const EventDetailLoading = () => (
  <Card className="event-detail">
    <p><span className="skeleton">Deskripsi event</span></p>
    <ul>
      <li><span className="skeleton">Waktu mulai: 01/01/2020 09:00</span></li>
      <li><span className="skeleton">Durasi: 0 jam 0 menit 0 detik</span></li>
    </ul>
  </Card>
);

const EventDetail = (props: EventDetailProps) => {
  const { event } = props;

  if (!event) return <EventDetailLoading />;

  const duration = durationText(event.startsAt, event.endsAt);
  return (
    <Card className="event-detail">
      <p>{event.description}</p>
      <ul>
        <li>Waktu mulai: {event.startsAt.toLocaleDateString()} {event.startsAt.toLocaleTimeString()}</li>
        <li>Durasi: {duration}</li>
      </ul>
    </Card>
  );
};

EventDetail.defaultProps = {
  event: null,
};

export default EventDetail;
