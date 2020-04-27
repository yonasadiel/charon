import React from 'react';
import { Card } from 'react-hephaestus';

const ParticipationLoadingPage = () => (
  <Card className="participation-page">
    <h1 className="title">Daftar Peserta</h1>
    <div className="participations">
      <div className="participation">
        <div className="user"><span className="skeleton">Nama peserta</span></div>
        <div className="venue"><span className="skeleton">Lokasi ujian</span></div>
      </div>
      <div className="participation">
        <div className="user"><span className="skeleton">Nama peserta</span></div>
        <div className="venue"><span className="skeleton">Lokasi ujian</span></div>
      </div>
      <div className="participation">
        <div className="user"><span className="skeleton">Nama peserta</span></div>
        <div className="venue"><span className="skeleton">Lokasi ujian</span></div>
      </div>
      <div className="participation">
        <div className="user"><span className="skeleton">Nama peserta</span></div>
        <div className="venue"><span className="skeleton">Lokasi ujian</span></div>
      </div>
    </div>
  </Card>
);

export default ParticipationLoadingPage;
