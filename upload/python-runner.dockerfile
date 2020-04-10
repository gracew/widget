FROM python-runner
ADD . .
RUN pip install -r requirements.txt