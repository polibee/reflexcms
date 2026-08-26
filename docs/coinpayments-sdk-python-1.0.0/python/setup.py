from setuptools import setup, find_packages

setup(
    name="coinpayments",
    version="1.0.0",
    packages=find_packages(),
    install_requires=["requests>=2.31"],
    python_requires=">=3.9",
)
