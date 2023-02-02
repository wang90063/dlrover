# Copyright 2023 The DLRover Authors. All rights reserved.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import unittest

import ray

from dlrover.python.scheduler.ray import RayClient
from dlrover.python.util.state.store_mananger import StoreManager


class RayClientTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        os.system("ray stop")
        r = os.system("ray start --head --port=5001  --dashboard-port=5000")
        assert r == 0
        ray.init("localhost:5001")
        cls.ray_client = RayClient.singleton_instance("test", "antc4mobius")
        cls.store_manager = StoreManager(
            jobname="antc4mobius", namespace="test"
        ).build_store_manager()
        cls.store = cls.store_manager.build_store()

    def test_create_and_delete_actor(self):
        class A:
            def __init__(self):
                self.a = "a"

            def run(self):
                return self.a

            def health_check(self, *args, **kargs):
                return "a"

        # to do: 使用类来描述启动的参数而不是dict

        actor_args = {"executor": A, "actor_name": "worker"}
        actor_handle = self.ray_client.create_actor(actor_args=actor_args)

        self.assertListEqual(self.store.get("actor_names"), ["worker"])
        res = getattr(actor_handle, "run").remote()
        self.assertEqual(ray.get(res), "a")

        for name, status in self.ray_client.list_actor():
            self.assertEqual(name, "worker")
            self.assertEqual(status, "RUNNING")

        self.ray_client.delete_actor(actor_args.get("actor_name"))
        self.assertListEqual(self.store.get("actor_names"), [])


if __name__ == "__main__":
    unittest.main()
