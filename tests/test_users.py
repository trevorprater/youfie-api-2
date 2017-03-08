import unittest
import requests
import json
from utils import utils


class TestCreateUserAndLoginLogout(unittest.TestCase):
    def setUp(self):
        utils.delete_user_if_exists('test1', 'venice')
        utils.delete_user_if_exists('test2', 'venice')

    def tearDown(self):
        utils.delete_user_if_exists('test1', 'venice')
        utils.delete_user_if_exists('test2', 'venice')

    def test_login(self):
        utils.create_user('test1', 'test1@youfie.io', 'venice')
        r, session = utils.login('test1', 'venice')
        self.assertEqual(r.status_code, 200)
        data = json.loads(r.content)
        self.assertTrue('token' in data.keys())
        self.assertTrue(len(data['token']) > 0)

    def test_logout(self):
        utils.create_user('test1', 'test1@youfie.io', 'venice')
        r, session = utils.login('test1', 'venice')
        self.assertEqual(r.status_code, 200)
        r = utils.view_user('test1', session)
        self.assertEqual(r.status_code, 200)
        utils.logout_user('test1', session)
        r = utils.view_user('test1', session)
        self.assertTrue(r.status_code != 200)

    def test_create_user(self):
        utils.create_user('test1', 'test1@youfie.io', 'venice')
        r, session = utils.login('test1', 'venice')
        self.assertEqual(r.status_code, 200)
        self.assertTrue('token' in json.loads(r.content).keys())

    def test_create_user_fails_duplicate_display_name(self):
        utils.create_user('test1', 'test1@youfie.io', 'venice')
        r = utils.create_user('test1', 'test7@youfie.io', 'venice')
        self.assertEqual(r.status_code, 409)

    def test_create_user_fails_duplicate_email(self):
        utils.create_user('test1', 'test1@youfie.io', 'venice')
        r = utils.create_user('test2', 'test1@youfie.io', 'venice')
        self.assertEqual(r.status_code, 409)

    def test_create_user_fails_no_password(self):
        r = utils.create_user('test1', 'test1@youfie.io', '')
        self.assertEqual(r.status_code, 500)


class TestUpdateUser(unittest.TestCase):
    def setUp(self):
        utils.delete_user_if_exists('test_update', 'venice')
        utils.delete_user_if_exists('trevorp', 'newpass')
        utils.create_user('test_update', 'test_update@youfie.io', 'venice')

    def tearDown(self):
        utils.delete_user_if_exists('test_update', 'venice')
        utils.delete_user_if_exists('trevorp', 'newpass')

    def test_update_user(self):
        r, session = utils.login('test_update', 'venice')
        initial_user = json.loads(
            utils.view_user('test_update', session).content)
        r = utils.update_user('test_update', {
            'password': 'newpass',
            'email': 'trevor.prater@gmail.com',
            'display_name': 'trevorp'
        }, session)
        user = json.loads(r.content)
        self.assertTrue(user['display_name'] != initial_user['display_name'])
        self.assertTrue(user['display_name'] == 'trevorp')
        self.assertTrue(user['email'] != initial_user['email'])

        r, session = utils.login('test_update', 'venice')
        self.assertEqual(r.status_code, 401)

        r, session = utils.login('trevorp', 'newpass')
        self.assertEqual(r.status_code, 200)


class TestViewUser(unittest.TestCase):
    def setUp(self):
        utils.delete_user_if_exists('test', 'venice')
        utils.delete_user_if_exists('test1', 'venice')
        utils.create_user('test', 'test@youfie.io', 'venice')
        utils.create_user('test1', 'test1@youfie.io', 'venice')

    def tearDown(self):
        utils.delete_user_if_exists('test', 'venice')
        utils.delete_user_if_exists('test1', 'venice')

    def test_view_self(self):
        r, session = utils.login('test', 'venice')
        r = utils.view_user('test', session)
        self.assertTrue(r.status_code, 200)
        data = json.loads(r.content)
        self.assertTrue(data['display_name'] == 'test')
        self.assertTrue(data['email'] == 'test@youfie.io')

    def test_view_other(self):
        r, session = utils.login('test', 'venice')
        r = utils.view_user('test1', session)
        self.assertEqual(r.status_code, 401)
