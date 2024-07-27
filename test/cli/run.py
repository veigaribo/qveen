import os
import shutil
import subprocess
import unittest  # Bad name
from typing import Callable, Type
from types import TracebackType


class TestDirContextManager:
	def __init__(self, test_name: str):
		self.test_name = test_name
		self.starting_cwd = os.getcwd()

	def __enter__(self):
		os.chdir(self.test_name)

	def __exit__(self,
							 exc_type: Type[BaseException] | None,
							 exc_value: BaseException | None,
							 exc_tb: TracebackType | None):
		os.chdir(self.starting_cwd)


def in_dir(dirname: str) -> TestDirContextManager:
	return TestDirContextManager(dirname)


def run_in_dir(dirname: str) -> Callable:
	def decorator(run: Callable) -> Callable:
		def wrapper(*args, **kwargs):
			with in_dir(dirname):
				return run(*args, **kwargs)

		return wrapper
	return decorator


class TestCli(unittest.TestCase):
	@run_in_dir('simple')
	def test_simple(self):
		subprocess.run(
			['qveen', 'params.toml'],
			check=True,
			stdout=subprocess.DEVNULL)

		with open('result.txt', encoding='utf-8') as file:
			result = file.read()
			self.assertEqual(result, '`something`\n')


if __name__ == '__main__':
	unittest.main()
