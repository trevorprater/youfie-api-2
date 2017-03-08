import unittest
import requests
import json

API_URL = 'http://localhost:5000'
TEST_USER = 'trevor'
TEST_USER_PW = 'venice'

def get_photos(display_name, session):
    return session.get(API_URL + '/users/{}/photos'.format(display_name))

def get_photo(display_name, photo_id, session):
    return session.get(API_URL + '/users/{}/photos/{}'.format(display_name, photo_id))

def create_photo(display_name, photo, session):
    return session.post(API_URL + '/users/{}/photos'.format(display_name), data=json.dumps(photo))

def delete_photo(display_name, photo_id, session):
    return session.delete(API_URL + '/users/{}/photos/{}'.format(display_name, photo_id))

def login(display_name, pw):
    r = requests.post(
        API_URL + '/users/{}/login'.format(display_name),
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
    r = requests.post(
        API_URL + '/users',
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
    return session.put(API_URL + '/users/{}'.format(display_name), data=json.dumps(updates))


def view_user(display_name, session):
    return session.get(API_URL + '/users/{}'.format(display_name))

def delete_user_if_exists(display_name, password):
    r, session = login(display_name, password)
    if r.status_code == 200:
        delete_user(display_name, session)
