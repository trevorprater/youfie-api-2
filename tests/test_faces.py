import unittest
import requests
import json
from utils import utils

class TestFace(unittest.TestCase):
    def setUp(self):
        utils.create_user('trevor', 'trevor@youfie.io', 'venice')
        _, self.session = utils.login('trevor', 'venice')
        self.photo = utils.PHOTO
        r = utils.create_photo('trevor', self.photo, self.session)
        self.photo['id'] = json.loads(r.content)['id']
        self.face = utils.FACE
        self.face['photo_id'] = self.photo['id']

    def tearDown(self):
        #utils.delete_user_if_exists('trevor', 'venice')
        utils.delete_photo('trevor', self.photo['id'], self.session)

    def test_create_face(self):
        r = utils.create_face('trevor', self.photo['id'], self.face, self.session)
        face = json.loads(r.content)
        self.assertEqual(r.status_code, 201)
        self.assertEqual(len(face['feature_vector']), 128)
        self.assertEqual(face['photo_id'], self.photo['id'])
        self.assertEqual(face['photo_id'], self.face['photo_id'])
        self.assertEqual(face['bb_top_left_x'], 120)
        self.assertEqual(face['bb_top_left_y'], 7)
        self.assertEqual(face['bb_top_right_x'], 256)
        self.assertEqual(face['bb_top_right_y'], 7)
        self.assertEqual(face['bb_bottom_left_x'], 120)
        self.assertEqual(face['bb_bottom_left_y'], 107)
        self.assertEqual(face['bb_bottom_right_x'], 256)
        self.assertEqual(face['bb_bottom_right_y'], 107)
