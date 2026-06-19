C Solution for Advent of Code 2015 day 8.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Part 1: sum of (code length - in-memory length). Part 2: sum of
C (encoded length - code length). Both are pure character counting; no
C string needs to be built, only measured.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER IOS, TOTAL, CODE, MEM, I
      CHARACTER*256 LINE
      CHARACTER C

      TOTAL = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
10    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 30
      CODE = LEN_TRIM(LINE)
      IF (CODE .EQ. 0) GOTO 10
C     Count in-memory characters, skipping the surrounding quotes.
      MEM = 0
      I = 2
20    IF (I .LE. CODE - 1) THEN
        C = LINE(I:I)
        IF (C .EQ. '\') THEN
          IF (LINE(I+1:I+1) .EQ. 'x') THEN
            I = I + 4
          ELSE
            I = I + 2
          END IF
        ELSE
          I = I + 1
        END IF
        MEM = MEM + 1
        GOTO 20
      END IF
      TOTAL = TOTAL + (CODE - MEM)
      GOTO 10
30    CLOSE(10)
      CALL ITOA(TOTAL, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER IOS, TOTAL, CODE, ENC, I
      CHARACTER*256 LINE
      CHARACTER C

      TOTAL = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
10    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 30
      CODE = LEN_TRIM(LINE)
      IF (CODE .EQ. 0) GOTO 10
C     Encoded length: 2 new surrounding quotes, plus each '"' or '\'
C     becomes two characters and every other character stays one.
      ENC = 2
      DO 20 I = 1, CODE
        C = LINE(I:I)
        IF (C .EQ. '"' .OR. C .EQ. '\') THEN
          ENC = ENC + 2
        ELSE
          ENC = ENC + 1
        END IF
20    CONTINUE
      TOTAL = TOTAL + (ENC - CODE)
      GOTO 10
30    CLOSE(10)
      CALL ITOA(TOTAL, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
