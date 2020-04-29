import React from 'react';
import { Card } from 'react-hephaestus';

const EventDetailLoading = () => (
  <Card className="event-detail">
    <p><span className="skeleton">Deskripsi event</span></p>
    <ul>
      <li><span className="skeleton">Waktu mulai: 01/01/2020 09:00</span></li>
      <li><span className="skeleton">Durasi: 0 jam 0 menit 0 detik</span></li>
    </ul>
  </Card>
);

export default EventDetailLoading;
