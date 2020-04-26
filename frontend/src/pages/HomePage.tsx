import React, { useEffect } from 'react';
import { Card } from 'react-hephaestus';
import { connect } from 'react-redux';

import CardMenu from '../components/navigation/CardMenu';
import conf from '../conf';
import { User, USER_ROLE } from '../modules/charon/auth/api';
import * as sessionSelectors from '../modules/session/selector';
import { AppState } from '../modules/store';
import { ROUTE_EXAM, ROUTE_VENUE, ROUTE_USER } from './routes';
import './HomePage.scss';

export interface ConnectedHomePageProps {
  user: User | null;
};

const EventMenu = () => <CardMenu to={ROUTE_EXAM} className="menu" iconName="fa-book" text="Ujian" />;
const UserMenu = () => <CardMenu to={ROUTE_USER} className="menu" iconName="fa-user-friends" text="Peserta" />;
const VenueMenu = () => <CardMenu to={ROUTE_VENUE} className="menu" iconName="fa-map-marker-alt" text="Lokasi Ujian" />;

const HomePage = (props: ConnectedHomePageProps) => {
  const { user } = props;
  useEffect(() => { document.title = conf.appName; }, []);

  return (
    <div className="home-page">
      <Card>
        <h1>{!!user ? `Selamat datang, ${user.name}` : conf.appName} </h1>
        {!!user
          ? (<div className="menus">
              <EventMenu />
              {user.role === USER_ROLE.ADMIN && <UserMenu />}
              {user.role === USER_ROLE.ADMIN && <VenueMenu />}
            </div>)
          : (<p>{conf.appName} adalah software untuk ujian semi-online.</p>)}
      </Card>
    </div>
  );
};

const mapStateToProps = (state: AppState) => ({
  user: sessionSelectors.getUser(state),
});

export default connect(mapStateToProps)(HomePage);
