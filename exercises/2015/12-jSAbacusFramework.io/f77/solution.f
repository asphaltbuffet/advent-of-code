C Solution for Advent of Code 2015 day 12.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Part 1: sum every number in the JSON document (a plain digit scan, since
C numbers only appear as values). Part 2: sum numbers but skip any object
C (and its descendants) that has a value of "red" — a recursive-descent
C walk with a shared cursor into the document buffer.

C     Shared document buffer and parse cursor.
      BLOCK DATA DAY12INIT
      INTEGER DLEN, POS
      CHARACTER DOC
      COMMON /D12/ DOC(40000), DLEN, POS
      END

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER DLEN, POS
      CHARACTER DOC
      COMMON /D12/ DOC(40000), DLEN, POS
      INTEGER TOTAL, I, NUM, SGN
      CHARACTER C
      LOGICAL INNUM

      CALL READDOC(INPUT_PATH)
      TOTAL = 0
      NUM = 0
      SGN = 1
      INNUM = .FALSE.
      DO 20 I = 1, DLEN
        C = DOC(I)
        IF (C .GE. '0' .AND. C .LE. '9') THEN
          NUM = NUM * 10 + (ICHAR(C) - ICHAR('0'))
          INNUM = .TRUE.
        ELSE
          IF (INNUM) THEN
            TOTAL = TOTAL + SGN * NUM
            NUM = 0
            INNUM = .FALSE.
          END IF
          IF (C .EQ. '-') THEN
            SGN = -1
          ELSE
            SGN = 1
          END IF
        END IF
20    CONTINUE
      IF (INNUM) TOTAL = TOTAL + SGN * NUM
      CALL ITOA(TOTAL, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER DLEN, POS
      CHARACTER DOC
      COMMON /D12/ DOC(40000), DLEN, POS
      INTEGER PARSEV

      CALL READDOC(INPUT_PATH)
      POS = 1
      CALL ITOA(PARSEV(), ANSWER, ANSWER_LEN)
      END

C     PARSEV — parse the value at the cursor, returning its number sum.
C     Objects containing a "red" value contribute 0; arrays always sum.
      RECURSIVE INTEGER FUNCTION PARSEV() RESULT(RES)
      INTEGER DLEN, POS
      CHARACTER DOC
      COMMON /D12/ DOC(40000), DLEN, POS
      CHARACTER C
      INTEGER SUM, SGN, NUM
      LOGICAL ISRED, ISREDV

      C = DOC(POS)
      IF (C .EQ. '{') THEN
        POS = POS + 1
        SUM = 0
        ISRED = .FALSE.
C       Object: comma-separated "key":value pairs until '}'.
10      IF (DOC(POS) .NE. '}') THEN
C         Key string.
          CALL SKIPSTR()
C         Colon.
          POS = POS + 1
C         Value: detect a literal "red" string value.
          IF (DOC(POS) .EQ. '"') THEN
            IF (ISREDV()) ISRED = .TRUE.
            CALL SKIPSTR()
          ELSE
            SUM = SUM + PARSEV()
          END IF
          IF (DOC(POS) .EQ. ',') POS = POS + 1
          GOTO 10
        END IF
        POS = POS + 1
        IF (ISRED) THEN
          RES = 0
        ELSE
          RES = SUM
        END IF
      ELSE IF (C .EQ. '[') THEN
        POS = POS + 1
        SUM = 0
20      IF (DOC(POS) .NE. ']') THEN
          IF (DOC(POS) .EQ. '"') THEN
            CALL SKIPSTR()
          ELSE
            SUM = SUM + PARSEV()
          END IF
          IF (DOC(POS) .EQ. ',') POS = POS + 1
          GOTO 20
        END IF
        POS = POS + 1
        RES = SUM
      ELSE IF (C .EQ. '"') THEN
        CALL SKIPSTR()
        RES = 0
      ELSE
C       A number (possibly negative).
        SGN = 1
        IF (C .EQ. '-') THEN
          SGN = -1
          POS = POS + 1
        END IF
        NUM = 0
30      C = DOC(POS)
        IF (C .GE. '0' .AND. C .LE. '9') THEN
          NUM = NUM * 10 + (ICHAR(C) - ICHAR('0'))
          POS = POS + 1
          GOTO 30
        END IF
        RES = SGN * NUM
      END IF
      END

C     ISREDV — true if the string starting at the cursor is exactly "red".
C     Does not advance the cursor.
      LOGICAL FUNCTION ISREDV()
      INTEGER DLEN, POS
      CHARACTER DOC
      COMMON /D12/ DOC(40000), DLEN, POS

      ISREDV = DOC(POS) .EQ. '"' .AND. DOC(POS+1) .EQ. 'r' .AND.
     &         DOC(POS+2) .EQ. 'e' .AND. DOC(POS+3) .EQ. 'd' .AND.
     &         DOC(POS+4) .EQ. '"'
      END

C     SKIPSTR — advance the cursor past a JSON string (cursor on opening ").
      SUBROUTINE SKIPSTR()
      INTEGER DLEN, POS
      CHARACTER DOC
      COMMON /D12/ DOC(40000), DLEN, POS

C     Skip the opening quote, then characters until the closing quote.
C     (AoC inputs contain no escaped quotes inside strings.)
      POS = POS + 1
10    IF (DOC(POS) .NE. '"') THEN
        POS = POS + 1
        GOTO 10
      END IF
      POS = POS + 1
      END

C     READDOC — load the single-line JSON document into the buffer.
      SUBROUTINE READDOC(INPUT_PATH)
      CHARACTER*(*) INPUT_PATH
      INTEGER DLEN, POS
      CHARACTER DOC
      COMMON /D12/ DOC(40000), DLEN, POS
      CHARACTER*40000 LINE
      INTEGER IOS, I

      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
      READ(10, '(A)', IOSTAT=IOS) LINE
      CLOSE(10)
      DLEN = LEN_TRIM(LINE)
      DO 10 I = 1, DLEN
        DOC(I) = LINE(I:I)
10    CONTINUE
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
