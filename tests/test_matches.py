import unittest
import requests
import json
from utils import utils


class TestMatch(unittest.TestCase):
    def setUp(self):
        utils.create_user('trevor', 'trevor@youfie.io', 'venice')
        _, self.session = utils.login('trevor', 'venice')
        self.user = json.loads(utils.get_user('trevor', self.session).content)
        self.photo = utils.PHOTO
        r = utils.create_photo('trevor', self.photo, self.session)
        photo = json.loads(r.content)
        self.photo['id'] = photo['id']
        self.face = utils.FACE
        self.face['photo_id'] = self.photo['id']
        r = utils.create_face('trevor', self.photo['id'], self.face,
                              self.session)
        face = json.loads(r.content)
        self.face['id'] = face['id']
        self.match = utils.MATCH
        self.match['user_id'] = self.user['id']
        self.match['photo_id'] = self.photo['id']
        self.match['face_id'] = self.face['id']

    def tearDown(self):
        utils.delete_face('trevor', self.photo['id'], self.face['id'],
                          self.session)
        utils.delete_photo('trevor', self.face['photo_id'], self.session)
        utils.delete_user_if_exists('trevor', 'venice')

    def test_get_match(self):
        r = utils.create_match('trevor', self.match, self.session)
        self.assertEqual(r.status_code, 201)
        match = json.loads(r.content)
        r = utils.get_match('trevor', match['id'], self.session)
        self.assertEqual(r.status_code, 200)
        recv_match = json.loads(r.content)
        self.assertEqual(match['id'], recv_match['id'])
        r = utils.delete_match('trevor', match['id'], self.session)
        self.assertEqual(r.status_code, 201)

    def test_get_matches(self):
        r = utils.get_matches('trevor', self.session)
        self.assertEqual(r.status_code, 200)
        recv_matches = json.loads(r.content)
        self.assertEqual(len(recv_matches), 0)
        created_match1 = json.loads(utils.create_match('trevor', self.match,
                                                       self.session).content)
        created_match2 = json.loads(utils.create_match('trevor', self.match,
                                                       self.session).content)
        potential_matches = json.loads(utils.get_potential_matches(
            'trevor', self.session).content)
        self.assertEqual(len(potential_matches), 2)
        recv_matches = json.loads(utils.get_matches('trevor',
                                                    self.session).content)
        self.assertEqual(len(recv_matches), 0)

        r = utils.delete_match('trevor', created_match1['id'], self.session)
        self.assertTrue(r.status_code, 201)
        r = utils.delete_match('trevor', created_match2['id'], self.session)
        self.assertTrue(r.status_code, 201)
        recv_matches = json.loads(utils.get_matches('trevor',
                                                    self.session).content)
        self.assertEqual(len(recv_matches), 0)
        potential_matches = json.loads(utils.get_potential_matches(
            'trevor', self.session).content)
        self.assertEqual(len(potential_matches), 0)

    def test_create_match(self):
        r = utils.create_match('trevor', self.match, self.session)
        self.assertEqual(r.status_code, 201)
        match = json.loads(r.content)
        r = utils.delete_match('trevor', match['id'], self.session)
        self.assertEqual(r.status_code, 201)

    def test_update_match(self):
        r = utils.create_match('trevor', self.match, self.session)
        self.assertEqual(r.status_code, 201)
        match = json.loads(r.content)
        match['is_match'] = False
        r = utils.update_match('trevor', match['id'], match, self.session)
        self.assertEqual(r.status_code, 201)
        r = utils.get_matches('trevor', self.session)
        self.assertEqual(r.status_code, 200)
        recv_matches = json.loads(r.content)
        self.assertEqual(len(recv_matches), 0)
        r = utils.delete_match('trevor', match['id'], self.session)
        self.assertEqual(r.status_code, 201)

    def test_delete_match(self):
        r = utils.create_match('trevor', self.match, self.session)
        self.assertEqual(r.status_code, 201)
        match = json.loads(r.content)
        r = utils.delete_match('trevor', match['id'], self.session)
        self.assertEqual(r.status_code, 201)
        recv_matches = json.loads(utils.get_matches('trevor',
                                                    self.session).content)
        self.assertEqual(len(recv_matches), 0)

    def test_viewing_matches(self):
        r = utils.create_match('trevor', self.match, self.session)
        self.assertEqual(r.status_code, 201)
        match = json.loads(r.content)

        r = utils.get_matches('trevor', self.session)
        matches = json.loads(r.content)
        self.assertEqual(len(matches), 0)

        r = utils.get_potential_matches('trevor', self.session)
        potential_matches = json.loads(r.content)
        self.assertEqual(len(potential_matches), 1)

        match['user_acknowledged'] = True
        r = utils.update_match('trevor', match['id'], match, self.session)
        self.assertEqual(r.status_code, 201)

        r = utils.get_potential_matches('trevor', self.session)
        potential_matches = json.loads(r.content)
        self.assertEqual(len(potential_matches), 0)

        r = utils.get_matches('trevor', self.session)
        matches = json.loads(r.content)
        self.assertEqual(len(matches), 1)

        r = utils.delete_match('trevor', match['id'], self.session)
        self.assertEqual(r.status_code, 201)
