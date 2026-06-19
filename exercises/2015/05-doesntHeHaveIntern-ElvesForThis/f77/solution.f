C Solution for Advent of Code 2015 day 5.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Count "nice" strings. Part 1: >=3 vowels, a doubled letter, and none of
C ab/cd/pq/xy. Part 2: a non-overlapping repeated pair, and a letter that
C repeats with exactly one letter between (xyx).

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER IOS, CNT
      CHARACTER*256 LINE
      LOGICAL NICE1

      CNT = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
10    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 20
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 10
      IF (NICE1(LINE, LEN_TRIM(LINE))) CNT = CNT + 1
      GOTO 10
20    CLOSE(10)
      CALL ITOA(CNT, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER IOS, CNT
      CHARACTER*256 LINE
      LOGICAL NICE2

      CNT = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
10    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 20
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 10
      IF (NICE2(LINE, LEN_TRIM(LINE))) CNT = CNT + 1
      GOTO 10
20    CLOSE(10)
      CALL ITOA(CNT, ANSWER, ANSWER_LEN)
      END

      LOGICAL FUNCTION NICE1(S, N)
      CHARACTER*256 S
      INTEGER N, I, NV
      CHARACTER A, B
      LOGICAL DBL, BAD

      NV = 0
      DBL = .FALSE.
      BAD = .FALSE.
      DO 10 I = 1, N
        A = S(I:I)
        IF (A.EQ.'a' .OR. A.EQ.'e' .OR. A.EQ.'i' .OR. A.EQ.'o'
     &      .OR. A.EQ.'u') NV = NV + 1
        IF (I .LT. N) THEN
          B = S(I+1:I+1)
          IF (A .EQ. B) DBL = .TRUE.
          IF ((A.EQ.'a'.AND.B.EQ.'b') .OR. (A.EQ.'c'.AND.B.EQ.'d')
     &        .OR. (A.EQ.'p'.AND.B.EQ.'q') .OR. (A.EQ.'x'.AND.B.EQ.'y'))
     &        BAD = .TRUE.
        END IF
10    CONTINUE
      NICE1 = (NV .GE. 3) .AND. DBL .AND. (.NOT. BAD)
      END

      LOGICAL FUNCTION NICE2(S, N)
      CHARACTER*256 S
      INTEGER N, I, J
      LOGICAL PAIR, REPEAT

      PAIR = .FALSE.
      REPEAT = .FALSE.
C     A repeated, non-overlapping pair: S(I:I+1) appears again at J >= I+2.
      DO 20 I = 1, N - 1
        DO 10 J = I + 2, N - 1
          IF (S(I:I+1) .EQ. S(J:J+1)) PAIR = .TRUE.
10      CONTINUE
20    CONTINUE
C     A letter that repeats with one between it (xyx).
      DO 30 I = 1, N - 2
        IF (S(I:I) .EQ. S(I+2:I+2)) REPEAT = .TRUE.
30    CONTINUE
      NICE2 = PAIR .AND. REPEAT
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
