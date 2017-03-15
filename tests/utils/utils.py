import unittest
import requests
import json
import random

API_URL = 'http://localhost:5000'
TEST_USER = 'trevor'
TEST_USER_PW = 'venice'

PHOTO = {
    'format': 'jpg',
    'storage_url': 'http://i.imgur.com/MsCzsxP.jpg',
    'latitude': 40.744381,
    'longitude': -73.987333,
    'width': 2448,
    'height': 3264,
}

FACE = {
    'feature_vector': json.dumps([round(
        random.uniform(0.1, 100.0), 6) for _ in xrange(128)]),
    'bb_top_left_x': 120,
    'bb_top_left_y': 7,
    'bb_top_right_x': 256,
    'bb_top_right_y': 7,
    'bb_bottom_left_x': 120,
    'bb_bottom_left_y': 107,
    'bb_bottom_right_x': 256,
    'bb_bottom_right_y': 107,
}

MATCH = {
    'photo_id': None,
    'face_id': None,
    'user_id': None,
    'confidence': 0.9,
    'user_acknowledged': False,
    'confirmed': False,
}


def create_match(display_name, match, session):
    return session.post(API_URL + '/users/{}/matches'.format(display_name),
                        data=json.dumps(match))


def get_matches(display_name, session):
    return session.get(API_URL + '/users/{}/matches'.format(display_name))


def get_match(display_name, match_id, session):
    return session.get(API_URL + '/users/{}/matches/{}'.format(display_name,
                                                               match_id))


def update_match(display_name, match_id, updates, session):
    return session.put(
        API_URL + '/users/{}/matches/{}'.format(display_name, match_id),
        data=json.dumps(updates))


def delete_match(display_name, match_id, session):
    return session.delete(API_URL + '/users/{}/matches/{}'.format(display_name,
                                                                  match_id))


def create_face(display_name, photo_id, face, session):
    return session.post(
        API_URL + '/users/{}/photos/{}/faces'.format(display_name, photo_id),
        data=json.dumps(face))


def get_faces(display_name, photo_id, session):
    return session.get(API_URL + '/users/{}/photos/{}/faces'.format(
        display_name, photo_id))


def get_face(display_name, photo_id, face_id, session):
    return session.get(API_URL + '/users/{}/photos/{}/faces/{}'.format(
        display_name, photo_id, face_id))


def update_face(display_name, photo_id, face_id, updates, session):
    return session.put(API_URL + '/users/{}/photos/{}/faces/{}'.format(
        display_name, photo_id, face_id),
                       data=json.dumps(updates))


def delete_face(display_name, photo_id, face_id, session):
    return session.delete(API_URL + '/users/{}/photos/{}/faces/{}'.format(
        display_name, photo_id, face_id))


def create_photo(display_name, photo, session):
    return session.post(API_URL + '/users/{}/photos'.format(display_name),
                        data=json.dumps(photo))


def get_photos(display_name, session):
    return session.get(API_URL + '/users/{}/photos'.format(display_name))


def get_photo(display_name, photo_id, session):
    return session.get(API_URL + '/users/{}/photos/{}'.format(display_name,
                                                              photo_id))


def update_photo(display_name, photo_id, updates, session):
    return session.put(
        API_URL + '/users/{}/photos/{}'.format(display_name, photo_id),
        data=json.dumps(updates))


def delete_photo(display_name, photo_id, session):
    return session.delete(API_URL + '/users/{}/photos/{}'.format(display_name,
                                                                 photo_id))


def login(display_name, pw):
    r = requests.post(API_URL + '/users/{}/login'.format(display_name),
                      data=json.dumps({
                          'password': pw
                      }))
    if r.status_code == 200:
        token = json.loads(r.content)['token']
        session = requests.Session()
        session.headers.update({'Authorization': 'Bearer {}'.format(token)})
        return r, session
    else:
        return r, None


def logout_user(display_name, session):
    return session.post(API_URL + '/users/{}/logout'.format(display_name))


def create_user(display_name, email, password):
    r = requests.post(API_URL + '/users',
                      data=json.dumps({
                          'password': password,
                          'display_name': display_name,
                          'email': email
                      }))
    return r


def delete_user(display_name, session):
    r = session.delete(API_URL + '/users/{}'.format(display_name))
    return r


def update_user(display_name, updates, session):
    return session.put(API_URL + '/users/{}'.format(display_name),
                       data=json.dumps(updates))


def get_user(display_name, session):
    return session.get(API_URL + '/users/{}'.format(display_name))


def delete_user_if_exists(display_name, password):
    r, session = login(display_name, password)
    if r.status_code == 200:
        delete_user(display_name, session)
